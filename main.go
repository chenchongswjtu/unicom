package main

import (
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	levelDB()
}

func levelDB() {
	//定义字符数组key value
	key := []byte("hello")
	value := []byte("hi i'm levelDB by go")
	//定义字符串的键值对
	k1 := "hhh"
	v1 := "heiheiheihei"

	// The returned DB instance is safe for concurrent use.
	// The DB must be closed after use, by calling Close method.
	//文件夹位置,获取db实际
	db, err := leveldb.OpenFile("/home/chenchong/gopath/src/awesomeProject/test/db", nil)

	//延迟关闭db流,必须的操作
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	//向db中插入键值对
	for i := 0; i < 1000; i++ {
		v := strconv.Itoa(i)
		e := db.Put(key, []byte(v), nil)
		//将字符串转换为byte数组
		e = db.Put([]byte(k1), []byte(v1), nil)
		fmt.Println(e)
	}
	e := db.Put(key, value, nil)
	//将字符串转换为byte数组
	e = db.Put([]byte(k1), []byte(v1), nil)
	fmt.Println(e) //<nil>
	/**
	  || 注意:查看路径下的db文件,可知道文件都是自动切分为一些小文件.
	*/
	//根据key值查询value
	data, _ := db.Get(key, nil)
	fmt.Println(data)        //[104 105 32 105 39 109 32 108 101 118 101 108 68 66 32 98 121 32 103 111]
	fmt.Printf("%c\n", data) //[h i   i ' m   l e v e l D B   b y   g o]
	data, _ = db.Get([]byte(k1), nil)
	fmt.Println(data)        //[104 101 105 104 101 105 104 101 105 104 101 105]
	fmt.Printf("%c\n", data) //[h e i h e i h e i h e i]
	// It is safe to modify the contents of the arguments after Delete returns but
	// not before.
	//根据key进行删除操作
	i := db.Delete(key, nil)
	fmt.Println(i) //<nil>
	data, err = db.Get(key, nil)
	fmt.Println(data, err) //[]
	//获取db快照
	snapshot, i := db.GetSnapshot()

	fmt.Println(snapshot.String()) //leveldb.Snapshot{22}

	newDB, err := leveldb.OpenFile("/home/chenchong/gopath/src/awesomeProject/test/newdb", nil)
	if err != nil {
		fmt.Println(err)
	}
	//延迟关闭db流,必须的操作
	defer newDB.Close()

	batch := new(leveldb.Batch)

	var res string
	iter := snapshot.NewIterator(nil, nil)
	for iter.Next() {
		res += fmt.Sprintf("(%s->%s)", string(iter.Key()), string(iter.Value()))
		batch.Put(iter.Key(), iter.Value())
	}

	iter.Release()

	fmt.Println(res)

	fmt.Println(i) //<nil>
	//注意: The snapshot must be released after use, by calling Release method.
	//也就是说snapshot在使用之后,必须使用它的Release方法释放!
	snapshot.Release()

	err = newDB.Write(batch, nil)
	if err != nil {
		panic(err)
		return
	}

	newSnapshot, err := newDB.GetSnapshot()
	if err != nil {
		panic(err)
		return
	}
	defer newSnapshot.Release()
	newIter := newSnapshot.NewIterator(nil, nil)
	defer newIter.Release()
	var newRes string
	for newIter.Next() {
		newRes += fmt.Sprintf("(%s->%s)", string(newIter.Key()), string(newIter.Value()))
	}

	fmt.Println(newRes)

}
