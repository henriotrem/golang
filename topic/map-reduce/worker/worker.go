package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"plugin"
	"sort"
	"strconv"
	"time"

	"github.com/henriotrem/topic/map-reduce/models"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrworker xxx.so\n")
		os.Exit(1)
	}

	mapf, reducef := loadPlugin(os.Args[1])

	Worker(mapf, reducef)
}

//
// load the application Map and Reduce functions
// from a plugin file, e.g. ../mrapps/wc.so
//
func loadPlugin(filename string) (func(string, string) []models.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []models.KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []models.KeyValue,
	reducef func(string, []string) string) {

	mr := models.MapReduce{
		Mapf:    mapf,
		Reducef: reducef,
	}

	for GetTask(mr) {
		time.Sleep(time.Second)
	}

}

func GetTask(mr models.MapReduce) bool {

	getTaskArgs := models.GetTaskArgs{}
	getTaskReply := models.GetTaskReply{}

	call("Master.GetTask", &getTaskArgs, &getTaskReply)

	task := getTaskReply.Task

	if task.Method == "map" {
		kva := ProcessInitialFile(mr, task)
		WriteIntermediateFiles(kva, task)
	} else if task.Method == "reduce" {
		intermediate := ProcessIntermediateFiles(task)
		WriteOutFile(mr, intermediate, task)
	} else if task.Method == "done" {
		return false
	} else {
		return true
	}

	updateTaskArgs := models.UpdateTaskArgs{}
	updateTaskReply := models.UpdateTaskReply{}

	updateTaskArgs.Task = getTaskReply.Task

	call("Master.UpdateTask", &updateTaskArgs, &updateTaskReply)

	return true
}

func ProcessInitialFile(mr models.MapReduce, task models.Task) []models.KeyValue {

	filename := task.Args[0]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mr.Mapf(filename, string(content))

	return kva
}

func WriteIntermediateFiles(kva []models.KeyValue, task models.Task) {

	nReduce, _ := strconv.Atoi(task.Args[1])
	mapTaskId := strconv.Itoa(task.Id)

	for i := 0; i < nReduce; i++ {
		reduceTaskId := strconv.Itoa(i)
		filename := "../data/intermediate/mr-" + mapTaskId + "-" + reduceTaskId
		tmp, err := ioutil.TempFile("../data/tmp", "mr-"+mapTaskId+"-"+reduceTaskId)
		if err != nil {
			fmt.Println(err)
			log.Fatalf("cannot create %v", filename)
		}
		enc := json.NewEncoder(tmp)
		for _, kv := range kva {
			if ihash(kv.Key)%nReduce == i {
				err := enc.Encode(&kv)
				if err != nil {
					log.Fatalf("cannot encode %v", &kv)
				}
			}
		}
		tmp.Close()
		if err := os.Rename(tmp.Name(), filename); err != nil {
			log.Printf("Rename error: %v", err)
			if i > 0 && !os.IsExist(err) {
				log.Println("Booo. IsExist() returned false, wanted true")
			}
		}
	}
}

func ProcessIntermediateFiles(task models.Task) []models.KeyValue {

	nMap, _ := strconv.Atoi(task.Args[0])
	reduceTaskId := strconv.Itoa(task.Id)
	intermediate := []models.KeyValue{}

	for i := 0; i < nMap; i++ {
		mapTaskId := strconv.Itoa(i)
		filename := "../data/intermediate/mr-" + mapTaskId + "-" + reduceTaskId
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		dec := json.NewDecoder(file)
		for {
			var kv models.KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
		file.Close()
	}

	sort.Sort(models.ByKey(intermediate))

	return intermediate
}

func WriteOutFile(mr models.MapReduce, intermediate []models.KeyValue, task models.Task) {

	oname := "../data/out/mr-out-" + strconv.Itoa(task.Id)
	ofile, _ := os.Create(oname)

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := mr.Reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := models.MasterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
