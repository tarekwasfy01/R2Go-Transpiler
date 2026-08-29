package runtime

import "fmt"

func init(){registerLoweringKernel("do_rowsum","0",kernelRowSum);registerLoweringKernel("do_rowsum","1",kernelRowSum)}
func kernelRowSum(c *Context,f *LoweringFrame)error{xv,e:=frameValue(c,f,0);if e!=nil{return e};dims,ok:=dimensions(xv);if !ok||len(dims)!=2{return fmt.Errorf("rowsum: matrix required")};x,e:=numbers(xv);if e!=nil{return e};gv,e:=frameValue(c,f,1);if e!=nil{return e};groups:=vectorKeys(gv);if len(groups)!=dims[0]{return fmt.Errorf("rowsum: group length mismatch")};order:=make([]string,0);index:=map[string]int{};for _,g:=range groups{if _,ok:=index[g];!ok{index[g]=len(order);order=append(order,g)}};out:=&DoubleVector{Data:make([]float64,len(order)*dims[1]),Missing:make([]bool,len(order)*dims[1])};for col:=0;col<dims[1];col++{for row,g:=range groups{i:=row+dims[0]*col;j:=index[g]+len(order)*col;if missingAt(x,i){out.Missing[j]=true;continue};out.Data[j]+=x.Data[i]}};_ = setAttribute(out,"dim",&IntegerVector{Data:[]int64{int64(len(order)),int64(dims[1])}});f.Result=out;return nil}
