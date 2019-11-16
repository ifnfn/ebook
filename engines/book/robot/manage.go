package robot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema/query"

	"ebook/engines/book"
	"ebook/resource/common"
)

// Library ...
type Library interface {
	Load()
	Save()
	Get(ID string) (book.Book, bool)
	Set(bok book.Book)
	Each(callback func(book book.Book) bool)
	Match(text string, offset, limit int) ([]book.Book, bool)
	Find(name string, value string) ([]book.Book, bool)
}

// LibraryDB ...
type LibraryDB struct {
	index   resource.Index
	books   *resource.Resource
	reviews *resource.Resource
}

// NewLibraryDB 新建书库
func NewLibraryDB(index resource.Index) Library {
	lib := &LibraryDB{
		index: index,
	}
	if books, exists := index.GetResource("books", nil); exists {
		lib.books = books
		lib.books.Compile(nil)
	} else {
		panic(exists)
	}

	if reviews, exists := index.GetResource("reviews", nil); exists {
		lib.reviews = reviews
		lib.reviews.Compile(nil)
	} else {
		panic(exists)
	}

	return lib
}

func bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return "[" + str + "]"
}

// Load ...
func (bm *LibraryDB) Load() {
	if data, err := ioutil.ReadFile("books.json"); err == nil {
		var books map[string]book.Book
		if json.Unmarshal(data, &books) == nil {
			ctx := context.TODO()
			items := []*resource.Item{}
			idx := 0
			size := len(books)

			for _, bok := range books {
				c := int(50 * float32(idx) / float32(size))
				fmt.Printf("\r%s %d%%", bar(c, 50), c*2)
				idx++

				payload := bok.Map()
				changes, base := bm.books.Validator().Prepare(ctx, payload, nil, false)
				doc, errs := bm.books.Validator().Validate(changes, base)
				if len(errs) == 0 {
					if item, err := resource.NewItem(doc); err == nil {
						items = append(items, item)
					}
				}
			}

			bm.books.Insert(ctx, items)
			fmt.Printf("\r%s %d%%\n", bar(50, 50), 100)
			fmt.Printf("共找到图书 %d 本\n", idx)
		}
	}
}

// Save ...
func (bm *LibraryDB) Save() {
	books := make(map[string]book.Book)
	bm.Each(func(bok book.Book) bool {
		books[bok.ID] = bok
		// system.PrintInterface(bok)
		return true
	})

	if data, err := json.MarshalIndent(books, "", "    "); err == nil {
		ioutil.WriteFile("books.json", data, 0666)
	} else {
		println(err.Error())
	}
}

// Get 根据 ID 查询图书
func (bm *LibraryDB) Get(ID string) (book.Book, bool) {
	bok := book.Book{}
	found := false

	if item, err := bm.books.Get(context.TODO(), ID); err == nil {
		bok.Fill(item.Payload)
		found = true
	}

	return bok, found
}

// Set ...
func (bm *LibraryDB) Set(bok book.Book) {
	q := &query.Query{}
	q.Predicate = append(q.Predicate, query.Equal{Field: "id", Value: bok.ID})

	payload := bok.Map()

	common.InsertOrPatchDB(context.TODO(), q, bm.books, payload)
}

// Each ...
func (bm *LibraryDB) Each(callback func(book book.Book) bool) {
	if items, err := bm.books.Find(context.TODO(), &query.Query{}); err == nil {
		for _, item := range items.Items {
			bok := book.Book{}
			bok.Fill(item.Payload)
			if callback(bok) == false {
				break
			}
		}
	}
}

// Find ...
func (bm *LibraryDB) Find(name string, value string) ([]book.Book, bool) {
	q := &query.Query{}
	q.Predicate = append(q.Predicate, query.Equal{Field: name, Value: value})

	// system.PrintInterface(query)
	books := make([]book.Book, 0)
	if items, err := bm.books.Find(context.TODO(), q); err == nil {
		for _, item := range items.Items {
			bok := book.Book{}
			bok.Fill(item.Payload)
			books = append(books, bok)
		}
	}

	return books, len(books) > 0
}

// Match 匹配书名
func (bm *LibraryDB) Match(text string, offset, limit int) ([]book.Book, bool) {
	q := &query.Query{}
	q.Predicate = append(q.Predicate, query.Regex{Field: "Name", Value: regexp.MustCompile(text)})
	q.Predicate = append(q.Predicate, query.Equal{Field: "Finish", Value: true})
	q.Window.Limit = limit
	q.Window.Offset = offset

	// query.SetSort()
	books := make([]book.Book, 0)
	if items, err := bm.books.Find(context.TODO(), q); err == nil {
		for _, item := range items.Items {
			bok := book.Book{}
			bok.Fill(item.Payload)
			books = append(books, bok)
		}
	}

	return books, len(books) > 0
}
