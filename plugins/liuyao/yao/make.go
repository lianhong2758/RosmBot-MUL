package yao

import (
	"math/rand/v2"
	"time"
)

//采用PCG随机数算法,用两个随机源作为标准

type Gua struct {
	// nowtime 调用时取即可
	// 用户数据的象征,作为其中一个种子
	Seed uint64
	//本卦
	//数字用来记录卦的阴阳,定义0为阳,1为阴
	//一爻为三数,最大为1+1+1,最小为0+0+0
	//则0为老阳,1为少阴,2为少阳.3为老阴
	Inward [6]uint8
	//变卦
	//老阳与老阴变换
	//结果再减去1,用于对应标准
	//最后应该只存在0,1
	//0对应 -- ,1对应 - -
	Range [6]uint8
	//变卦数量
	RangeNum int
}
type Pair struct {
	A uint8
	B uint8
}

func NewGua(seed uint64) *Gua {
	return &Gua{
		Seed:   seed,
		Inward: [6]uint8{0, 0, 0, 0, 0, 0},
		Range:  [6]uint8{0, 0, 0, 0, 0, 0},
	}
}

// 起卦
func (g *Gua) Divination(id int) {
	var x uint8 = 0
	for range 3 {
		r := rand.New(rand.NewPCG(g.Seed, uint64(time.Now().UnixNano())))
		x += uint8(r.UintN(2))                                     //0/1
		time.Sleep(time.Nanosecond * time.Duration(rand.IntN(20))) //随机暂停时间用来避免干扰性
	}
	g.Inward[6-1-id] = x
	return
}

// 变卦
func (g *Gua) Changes() {
	g.RangeNum = 0
	for k, v := range g.Inward {
		switch v {
		case 0:
			g.Range[k] = 1
			g.RangeNum++
		case 3:
			g.Range[k] = 0
			g.RangeNum++
		case 1: //少阴
			g.Range[k] = 1
		case 2: //少阳
			g.Range[k] = 0
		}
	}
}

func (g *Gua) GetManifestation() Pair {
	return Pair{
		A: g.Range[0]<<2 + g.Range[1]<<1 + g.Range[2],
		B: g.Range[3]<<2 + g.Range[4]<<1 + g.Range[5],
	}
}

//获取本卦输出
func (g *Gua) GetInward ()  string {
	s:=""
	 for _,v:=range g.Inward{
		 switch v{
		 case 0:
			 s+="太阳"
		 case 1:
			 s+="少阴"
		 case 2:
			 s+="少阳"
		 case 3:
			 s+="太阴"
		 }
	 }
	 return s
}
//获取变卦输出
func (g *Gua) GetRange()  string {
	s:=""
	 for _,v:=range g.Range{
		s+=GetName(v)
	 }
	 return s
}