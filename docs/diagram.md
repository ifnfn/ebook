# Title
<!-- TOC -->

- [Title](#title)
    - [aaa](#aaa)
        - [CC](#cc)
- [](#)
- [](#-1)

<!-- /TOC -->

## aaa
### CC

>$${\bold e}_\alpha = \begin{pmatrix}\cos\alpha\\ \sin\alpha\end{pmatrix}$$  

##

$$a{\bold e}_\alpha + b{\bold e}_\alpha + ce_x=0$$
##

 > $$e^2 = \sqrt{a^2+b^2}$$ (2)

```mermaid
graph TB;
    A[Hard edge] -->|Link text| B(Round edge)
    B --> C{Decision}
    C -->|Y| B
    C -->|N| A  
```

```mermaid
sequenceDiagram
    participant H as 防丢片
    participant M as 手机
    participant S as 服务器
    M->>H: 扫描蓝牙
    H->>M: service
    M->>H: 连接
    M->>S: 检查防丢片是否注册
    M->>S: 登记防丢片(基本信息)
    Note right of S: ?
```

```mermaid
sequenceDiagram
    participant H as 防丢片
    participant M as 手机
    participant S as 服务器
    loop 定期检测
        M->>S: 读取我的设备列表
        S->>M: 更新本地设备列表
        M->>H: 扫描蓝牙
        H->>M: service
        M->>H: 连接
        H->>M: 上报信息
        M->>M: 蓝牙设备检控
        M->>S: 更新防丢片状态
        M->>M: 防丢片失联处理
    end
```

```mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Alice->John: Hello John, how are you?
    loop Healthcheck
        John->John: Fight against hypochondria
    end
    Note right of John: Rational thoughts <br/>prevail...
    John-->Alice: Great!
    John->Bob: How about you?
    Bob-->John: Jolly good!
```

```mermaid
gantt
        dateFormat  YYYY-MM-DD
        title Adding GANTT diagram functionality to mermaid
        section A section
        Completed task            :done,    des1, 2014-01-06,2014-01-08
        Active task               :active,  des2, 2014-01-09, 3d
        Future task               :         des3, after des2, 5d
        Future task2               :         des4, after des3, 5d
        section Critical tasks
        Completed task in the critical line :crit, done, 2014-01-06,24h
        Implement parser and jison          :crit, done, after des1, 2d
        Create tests for parser             :crit, active, 3d
        Future task in critical line        :crit, 5d
        Create tests for renderer           :2d
        Add to mermaid                      :1d
```

```mermaid
sequenceDiagram
  participant  A as Client
  participant  B as Server
  participant  C as Goods
  A->B: Works!
  A->>B: Works!
  A-->B: Works!
  A-->>B: Works!
  B->>B: Works!
  A->>B: Works!
  A->>A: Works!
  A->>B: Works!
  B->>C: Query
  C->>A: Return
  Note left of A: Note left
  Note over A,C: Over message infomation
  Note right of C: Thinking ....
```

```mermaid
graph
A->B;
B->C;
```

```mermaid
sequenceDiagram
Andrew->China: Says Hello 
Note right of China: China thinks\nabout it 
China-->Andrew: How are you? 
Andrew->>China: I am good thanks!
```