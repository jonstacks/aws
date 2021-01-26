package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dustin/go-humanize"
	"github.com/jonstacks/aws/pkg/models"
	"github.com/jonstacks/aws/pkg/utils"
)

type VFStat struct {
	count uint64
	size  uint64
}

type VirtualFolder struct {
	objects  *VFStat
	Prefixes map[string]*VFStat
}

type VFStatPair struct {
	Key   string
	Value *VFStat
}
type VFStatPairList []VFStatPair

func (p VFStatPairList) Len() int           { return len(p) }
func (p VFStatPairList) Less(i, j int) bool { return p[i].Value.size < p[j].Value.size }
func (p VFStatPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (vf *VirtualFolder) AddObject(o *s3.Object) {
	vf.objects.count++
	vf.objects.size += uint64(aws.Int64Value(o.Size))
}

func (vf *VirtualFolder) AddObjectWithPrefix(prefix string, o *s3.Object) {
	_, ok := vf.Prefixes[prefix]
	if !ok {
		vf.Prefixes[prefix] = &VFStat{}
	}
	vf.Prefixes[prefix].count++
	vf.Prefixes[prefix].size += uint64(aws.Int64Value(o.Size))
}

func (vf *VirtualFolder) TotalSize() uint64 {
	var total uint64
	total += vf.objects.size
	for _, group := range vf.Prefixes {
		total += group.size
	}
	return total
}

func (vf *VirtualFolder) TotalObjects() uint64 {
	var count uint64
	count += vf.objects.count
	for _, group := range vf.Prefixes {
		count += group.count
	}
	return count
}

func main() {
	models.Init(models.DefaultSession())

	bucket := os.Args[1]
	prefix := ""
	if len(os.Args) > 2 {
		prefix = os.Args[2]
	}

	vf := &VirtualFolder{
		objects:  &VFStat{},
		Prefixes: make(map[string]*VFStat),
	}

	objChan := make(chan *s3.Object, 100000)
	done := false
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Printf("%d objects metadata scanned (%s)...\n", vf.TotalObjects(), humanize.Bytes(vf.TotalSize()))
			if done {
				return
			}
		}
	}()

	go func() {
		err := models.GetS3ObjectsWithPrefixChan(bucket, prefix, objChan)
		if err != nil {
			utils.ExitErrorHandler(err)
		}
	}()

	for object := range objChan {
		key := aws.StringValue(object.Key)
		relPrefix := strings.TrimPrefix(key, fmt.Sprintf("%s/", prefix))

		if strings.Contains(relPrefix, "/") {
			vf.AddObjectWithPrefix(strings.Split(relPrefix, "/")[0], object)
		} else {
			vf.AddObject(object)
		}
	}

	fmt.Printf("s3://%s/%s\n", bucket, prefix)
	fmt.Printf("--------------------------\n")

	if len(vf.Prefixes) > 0 {
		fmt.Printf("Prefixes:\n")
		orderedBySize := make(VFStatPairList, len(vf.Prefixes))
		i := 0
		for prefix, group := range vf.Prefixes {
			orderedBySize[i] = VFStatPair{Key: prefix, Value: group}
			i++
		}
		sort.Sort(orderedBySize)
		for _, vfpair := range orderedBySize {
			fmt.Printf("%s\t%s\n", humanize.Bytes(vfpair.Value.size), vfpair.Key)
		}
		fmt.Printf("--------------------------\n")
	}

	fmt.Printf("Total Count: %d\n", vf.TotalObjects())
	fmt.Printf("Total Size: %s\n", humanize.Bytes(vf.TotalSize()))
}
