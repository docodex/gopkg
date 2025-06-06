# gopkg

`gopkg` is a universal utility collection for Go.

## Table of Contents

- [Introduction](#Introduction)
- [Catalogs](#Catalogs)
- [How To Use](#How-To-Use)
- [License](#License)

## Introduction

`gopkg` is a universal utility collection for Go, it complements offerings such as Better std, Data Structures, and so on.

## Catalogs

- [Container](https://github.com/docodex/gopkg/tree/master/container): Data Structures
    - [List](https://github.com/docodex/gopkg/tree/master/container/list)
        - [ArrayList](https://github.com/docodex/gopkg/tree/master/container/list/arraylist): array list
        - [DoublyLinkedList](https://github.com/docodex/gopkg/tree/master/container/list/doublylinkedlist): doubly linked list
        - [SinglyLinkedList](https://github.com/docodex/gopkg/tree/master/container/list/singlylinkedlist): singly linked list
    - [Set](https://github.com/docodex/gopkg/tree/master/container/set)
        - [HashSet](https://github.com/docodex/gopkg/tree/master/container/set/hashset): set backed by hash table
        - [TreeSet](https://github.com/docodex/gopkg/tree/master/container/set/treeset): set backed by red-black tree
    - [Map](https://github.com/docodex/gopkg/tree/master/container/dict)
        - [HashMap](https://github.com/docodex/gopkg/tree/master/container/dict/hashmap): map backed by hash table
        - [TreeMap](https://github.com/docodex/gopkg/tree/master/container/dict/treemap): map backed by red-black tree
        - [HashBidiMap](https://github.com/docodex/gopkg/tree/master/container/dict/hashbidimap): bidirectional map backed by hash tables
        - [TreeBidiMap](https://github.com/docodex/gopkg/tree/master/container/dict/treebidimap): bidirectional map backed by red-black trees
    - [Stack](https://github.com/docodex/gopkg/tree/master/container/stack)
        - [ArrayStack](https://github.com/docodex/gopkg/tree/master/container/stack/arraystack): array stack
        - [LinkedListStack](https://github.com/docodex/gopkg/tree/master/container/stack/linkedliststack): singly linked list stack
    - [Queue](https://github.com/docodex/gopkg/tree/master/container/queue)
        - [ArrayQueue](https://github.com/docodex/gopkg/tree/master/container/queue/arrayqueue): array queue
        - [LinkedListQueue](https://github.com/docodex/gopkg/tree/master/container/queue/linkedlistqueue): singly linked list queue
        - [DoubleEndedQueue](https://github.com/docodex/gopkg/tree/master/container/queue/deque): double ended queue
        - [CircularQueue](https://github.com/docodex/gopkg/tree/master/container/queue/circularqueue): circular buffer (circular queue)
        - [PriorityQueue](https://github.com/docodex/gopkg/tree/master/container/queue/priorityqueue): priority queue (binary heap)
    - [Tree](https://github.com/docodex/gopkg/tree/master/container/tree)
        - [AVLTree](https://github.com/docodex/gopkg/tree/master/container/tree/avltree): AVL balanced binary tree
        - [RedBlackTree](https://github.com/docodex/gopkg/tree/master/container/tree/redblacktree): red-black tree
        - [BTree](https://github.com/docodex/gopkg/tree/master/container/tree/btree): B-tree
        - [BinaryHeap](https://github.com/docodex/gopkg/tree/master/container/tree/binaryheap): binary heap
    - [Ring](https://github.com/docodex/gopkg/tree/master/container/ring)
        - [DoublyLinkedRing](https://github.com/docodex/gopkg/tree/master/container/ring/doublylinkedring): doubly linked circular list
        - [SinglyLinkedRing](https://github.com/docodex/gopkg/tree/master/container/ring/singlylinkedring): singly linked circular list
    - [Skiplist](https://github.com/docodex/gopkg/tree/master/container/skiplist): skip list (skiplist)
- [Snowflake](https://github.com/docodex/gopkg/tree/master/snowflake): ID Generator in Twitter snowflake format

## How To Use

You can use `go get -u github.com/docodex/gopkg@master` to get or update `gopkg`.

## License

`gopkg` is licensed under the terms of the MIT license. See [LICENSE](LICENSE) for more information.
