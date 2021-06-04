/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ast

import (
    "strconv"
    "testing"

    "github.com/bytedance/sonic/internal/native"
    jsoniter "github.com/json-iterator/go"
)

var parallelism = 4

func TestNodeRaw(t *testing.T) {
    root, derr := NewSearcher(_TwitterJson).GetByPath("search_metadata")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val := root.Raw()
    var comp = `{
    "max_id": 250126199840518145,
    "since_id": 24012619984051000,
    "refresh_url": "?since_id=250126199840518145&q=%23freebandnames&result_type=mixed&include_entities=1",
    "next_results": "?max_id=249279667666817023&q=%23freebandnames&count=4&include_entities=1&result_type=mixed",
    "count": 4,
    "completed_in": 0.035,
    "since_id_str": "24012619984051000",
    "query": "%23freebandnames",
    "max_id_str": "250126199840518145"
  }`
    if val != comp {
        t.Fatalf("exp: %+v, got: %+v", comp, val)
    }

    root, derr = NewSearcher(_TwitterJson).GetByPath("statuses", 0, "entities", "hashtags")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val = root.Raw()
    comp = `[
          {
            "text": "freebandnames",
            "indices": [
              20,
              34
            ]
          }
        ]`
    if val != comp {
        t.Fatalf("exp: \n%s\n, got: \n%s\n", comp, val)
    }

    var data = `{"k1" :   { "a" : "bc"} , "k2" : [1,2 ] , "k3":{} , "k4":[]}`
    root, derr = NewSearcher(data).GetByPath("k1")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val = root.Raw()
    comp = `{ "a" : "bc"}`
    if val != comp {
        t.Fatalf("exp: %+v, got: %+v", comp, val)
    }

    root, derr = NewSearcher(data).GetByPath("k2")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val = root.Raw()
    comp = `[1,2 ]`
    if val != comp {
        t.Fatalf("exp: %+v, got: %+v", comp, val)
    }

    root, derr = NewSearcher(data).GetByPath("k3")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val = root.Raw()
    comp = `{}`
    if val != comp {
        t.Fatalf("exp: %+v, got: %+v", comp, val)
    }

    root, derr = NewSearcher(data).GetByPath("k4")
    if derr != nil {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val = root.Raw()
    comp = `[]`
    if val != comp {
        t.Fatalf("exp: %+v, got: %+v", comp, val)
    }
}

func TestNodeGet(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val := root.Get("search_metadata").Get("max_id").Int64()
    if val != int64(250126199840518145) {
        t.Fatalf("exp: %+v, got: %+v", 250126199840518145, val)
    }
}

func TestNodeIndex(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val := root.Get("statuses").Index(3).Get("id_str").String()
    if val != "249279667666817024" {
        t.Fatalf("exp: %+v, got: %+v", "249279667666817024", val)
    }
}

func TestNodeGetByPath(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    val := root.GetByPath("statuses", 3, "id_str").String()
    if val != "249279667666817024" {
        t.Fatalf("exp: %+v, got: %+v", "249279667666817024", val)
    }
}

func TestNodeSet(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    root.GetByPath("statuses", 3).Set("id_str", Node{
        v: int64(111),
        t: native.V_INTEGER,
    })
    val := root.GetByPath("statuses", 3, "id_str").Int64()
    if val != 111 {
        t.Fatalf("exp: %+v, got: %+v", 111, val)
    }
    for i := root.GetByPath("statuses", 3).Cap(); i >= 0; i-- {
        root.GetByPath("statuses", 3).Set("id_str"+strconv.Itoa(i), Node{
            v: int64(111),
            t: native.V_INTEGER,
        })
    }
    val = root.GetByPath("statuses", 3, "id_str0").Int64()
    if val != 111 {
        t.Fatalf("exp: %+v, got: %+v", 111, val)
    }

    nroot, derr := NewParser(`{"a":[0.1,true,0,"name",{"b":"c"}]}`).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    root.GetByPath("statuses", 3).Set("id_str2", nroot)
    val2 := root.GetByPath("statuses", 3, "id_str2", "a", 4, "b").String()
    if val2 != "c" {
        t.Fatalf("exp:%+v, got:%+v", "c", val2)
    }
}

func TestNodeSetByIndex(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    root.GetByPath("statuses").SetByIndex(0, Node{
        v: int64(111),
        t: native.V_INTEGER,
    })
    val := root.GetByPath("statuses", 0).Int64()
    if val != 111 {
        t.Fatalf("exp: %+v, got: %+v", 111, val)
    }

    nroot, derr := NewParser(`{"a":[0.1,true,0,"name",{"b":"c"}]}`).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    root.GetByPath("statuses").SetByIndex(0, nroot)
    val2 := root.GetByPath("statuses", 0, "a", 4, "b").String()
    if val2 != "c" {
        t.Fatalf("exp:%+v, got:%+v", "c", val2)
    }
}

func TestNodeAdd(t *testing.T) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }

    for i := root.GetByPath("statuses").Cap(); i >= 0; i-- {
        root.GetByPath("statuses").Add(Node{
            v: int64(111),
            t: native.V_INTEGER,
        })
    }
    val := root.GetByPath("statuses", 4).Int64()
    if val != 111 {
        t.Fatalf("exp: %+v, got: %+v", 111, val)
    }
    val = root.GetByPath("statuses", root.GetByPath("statuses").Len()-1).Int64()
    if val != 111 {
        t.Fatalf("exp: %+v, got: %+v", 111, val)
    }

    nroot, derr := NewParser(`{"a":[0.1,true,0,"name",{"b":"c"}]}`).Parse()
    if derr != 0 {
        t.Fatalf("decode failed: %v", derr.Error())
    }
    root.GetByPath("statuses").Add(nroot)
    val2 := root.GetByPath("statuses", root.GetByPath("statuses").Len()-1, "a", 4, "b").String()
    if val2 != "c" {
        t.Fatalf("exp:%+v, got:%+v", "c", val2)
    }
}

func BenchmarkNodeRaw(b *testing.B) {
    root, derr := NewSearcher(_TwitterJson).GetByPath("search_metadata")
    if derr != nil {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    b.SetParallelism(parallelism)
    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            root.Raw()
        }
    })
}

func BenchmarkNodeGetByPath(b *testing.B) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    _ = root.GetByPath("statuses", 3, "entities", "hashtags", 0, "text").String()
    b.SetParallelism(parallelism)
    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _ = root.GetByPath("statuses", 3, "entities", "hashtags", 0, "text").String()
        }
    })
}

func BenchmarkStructGetByPath_Jsoniter(b *testing.B) {
    var root = _TwitterStruct{}
    err := jsoniter.Unmarshal([]byte(_TwitterJson), &root)
    if err != nil {
        b.Fatalf("unmarshal failed: %v", err)
    }

    b.SetParallelism(parallelism)
    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _ = root.Statuses[3].Entities.Hashtags[0].Text
        }
    })
}

func BenchmarkNodeGet(b *testing.B) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    node := root.Get("statuses").Index(3).Get("entities").Get("hashtags").Index(0)
    node.Set("test1", newInt64(1))
    node.Set("test2", newInt64(2))
    node.Set("test3", newInt64(3))
    node.Set("test4", newInt64(4))
    node.Set("test5", newInt64(5))
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        node.Get("text").Interface()
    }
}

func BenchmarkMapGet(b *testing.B) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    node := root.Get("statuses").Index(3).Get("entities").Get("hashtags").Index(0)
    node.Set("test1", newInt64(1))
    node.Set("test2", newInt64(2))
    node.Set("test3", newInt64(3))
    node.Set("test4", newInt64(4))
    node.Set("test5", newInt64(5))
    m := node.Map()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m["text"]
    }
}

func BenchmarkNodeSet(b *testing.B) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    node := root.Get("statuses").Index(3).Get("entities").Get("hashtags").Index(0)
    b.SetParallelism(parallelism)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        node.Set("test1", newInt64(1))
    }
}

func BenchmarkMapSet(b *testing.B) {
    root, derr := NewParser(_TwitterJson).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    node := root.Get("statuses").Index(3).Get("entities").Get("hashtags").Index(0).Map()
    b.SetParallelism(parallelism)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        node["test1"] = map[string]int{"test1": 1}
    }
}

func BenchmarkNodeAdd(b *testing.B) {
    data := `{"statuses":[]}`
    _, derr := NewParser(data).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    b.SetParallelism(parallelism)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        root, _ := NewParser(data).Parse()
        node := root.Get("statuses")
        node.Add(newObject([]Pair{{"test", newInt64(1)}}))
    }
}

func BenchmarkSliceAdd(b *testing.B) {
    data := `{"statuses":[]}`
    _, derr := NewParser(data).Parse()
    if derr != 0 {
        b.Fatalf("decode failed: %v", derr.Error())
    }
    b.SetParallelism(parallelism)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        root, _ := NewParser(data).Parse()
        node := root.Get("statuses").Array()
        _ = append(node, map[string]interface{}{"test": 1})
    }
}