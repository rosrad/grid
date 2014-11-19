package monc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"util"
)

func GenerateList(inlist, outlist string) bool {
	if !util.IsExist(inlist) {
		fmt.Println("input list do not exist! File : ", inlist)
		return false
	}
	ofw, err := os.Create(outlist)
	if err != nil {
		fmt.Println("output list generatation failed! File: ", outlist)
		return false
	}
	defer ofw.Close()
	ofbw := bufio.NewWriter(ofw)
	ifr, err := os.Open(inlist)
	if err != nil {
		fmt.Println("input list can not open! File: ", inlist)
		return false
	}
	defer ifr.Close()
	ifbr := bufio.NewScanner(ifr)
	for ifbr.Scan() {
		key := ifbr.Text()
		fmt.Println("Keyword: ", key)
		str := fmt.Sprintf("/CDShare/Corpus/monc/clean/NU-%s.raw\n", key)
		fmt.Println("Path: ", str)
		ofbw.WriteString(str)
	}
	ofbw.Flush()
	return true
}

func FileSize(file string) int64 {
	if !util.IsExist(file) {
		fmt.Println("No Exist File : ", file)
		return -1
	}
	info, err := os.Stat(file)
	if err != nil {
		fmt.Println("Info no catched File : ", file)
		return -1
	}
	return info.Size()
}

func SLen(size int64, rate int, b int) float64 {
	return float64(size) / float64(rate*b/8)
}

func MSize(time float64, rate int, b int) int64 {
	sz:= int64(float64(rate*b/8) * time)
	return (sz/2+sz%2)*2
}
func MergeRaw(inlist, outdir string) bool {
	if !util.IsExist(inlist) {
		return false
	}
	ifr, err := os.Open(inlist)
	if err != nil {
		return false
	}
	defer ifr.Close()
	var maxSize int64 = MSize(15*60, 8000, 16)
	var sumSize int64 = 0
	var index int = 0
	var fmerge, flist *os.File = nil, nil
	var fw, flw *bufio.Writer = nil, nil
	var zerobuf [8000 * 2 / 4]byte
	isc := bufio.NewScanner(ifr)
	for isc.Scan() {
		fileSize := FileSize(isc.Text())
		if sumSize <= 0 || sumSize+fileSize > maxSize {
			if fmerge != nil {
				if fw != nil {
					fw.Flush()
					flw.Flush()
				} else {
					fmt.Println(" bufio Write Error")
				}
				flist.Close()
				fmerge.Close()
			}
			// update merge file
			mergefile := filepath.Join(outdir, fmt.Sprintf("merge%d.raw", index))
			filelist := filepath.Join(outdir, fmt.Sprintf("merge%d.list", index))
			fmt.Println("merge file ", mergefile)
			flist, _ = os.Create(filelist)
			fmerge, _ = os.Create(mergefile)
			fw = bufio.NewWriter(fmerge)
			flw = bufio.NewWriter(flist)
			index++
			sumSize = 0
		}
		fmt.Println("File: ", isc.Text(), "\n Size : ", fileSize)
		list := fmt.Sprintf("%s %f\n", isc.Text(), SLen(fileSize, 8000, 16))
		fr, _ := os.Open(isc.Text())
		io.Copy(fw, fr)
		fw.Write(zerobuf[0:])
		flw.WriteString(list)
		sumSize = sumSize + fileSize
	}
	if fmerge != nil {
		fmt.Println("the last file close")

		if fw != nil {
			fw.Flush()
			flw.Flush()
		} else {
			fmt.Println(" bufio Write Error")
		}
		fmerge.Close()
		flist.Close()
	}
	return true
}

func CutRaw(list, mer_file, outdir string) bool {
	f_list, err := os.Open(list)
	if err != nil {
		fmt.Println("Can not open the list ,err: ", err)
	}
	defer f_list.Close()
	// read input contracted sound file
	f_mer, err := os.Open(mer_file)
	if err != nil {
		fmt.Println("Can not open the merged file, err: ", err)
	}
	defer f_mer.Close()

	scanner := bufio.NewScanner(f_list)
	//skip_sz := MSize(0.25, 8000, 16)
	//skip_buf := make([]byte, skip_sz)
	
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		out_f := filepath.Join(outdir, filepath.Base(fields[0]))
		f_sz, _ := strconv.ParseFloat(fields[1], 10)
		//fmt.Printf("Size %f ,File: %s\n", f_sz, out_f)
		ofw, err_w := os.Create(out_f)
		if err_w != nil {
			fmt.Println("Can not Create output file, err: ", err_w)
		}
		defer ofw.Close()
		
		wr_size := MSize(f_sz+0.25, 16000, 16)
		f_lim := io.LimitReader(f_mer, wr_size)
		wrt_n, err_cp:=io.Copy(ofw, f_lim)
		fmt.Printf("Write File : %s\n Need : %d , Written %d, Align %d\n",out_f, wr_size, wrt_n, wrt_n%2)
		if err_cp !=nil || wr_size != wrt_n{
			fmt.Println("Write to the splited file failed : err", err_cp)
			fmt.Printf("Need to Write %d, Writed %d", wr_size,wrt_n)
		}
		//f_skip := io.LimitReader(f_mer, skip_sz)
		//f_skip.Read(skip_buf)
	}
	return true
}




