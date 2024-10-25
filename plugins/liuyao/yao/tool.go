package yao

var GuaMap = map[Pair]int{
	{writes[0], writes[0]}: 0,
	{writes[7], writes[7]}: 1,
	{writes[3], writes[5]}: 2,
	{writes[5], writes[6]}: 3,
	{writes[0], writes[5]}: 4,
	{writes[5], writes[0]}: 5,
	{writes[5], writes[7]}: 6,
	{writes[7], writes[5]}: 7,
	{writes[0], writes[4]}: 8,
	{writes[1], writes[0]}: 9,
	{writes[0], writes[7]}: 10,
	{writes[7], writes[0]}: 11,
	{writes[2], writes[0]}: 12,
	{writes[0], writes[2]}: 13,
	{writes[6], writes[7]}: 14,
	{writes[7], writes[3]}: 15,
	{writes[3], writes[1]}: 16,
	{writes[4], writes[6]}: 17,
	{writes[1], writes[7]}: 18,
	{writes[7], writes[4]}: 19,
	{writes[3], writes[2]}: 20,
	{writes[2], writes[6]}: 21,
	{writes[7], writes[6]}: 22,
	{writes[3], writes[7]}: 23,
	{writes[3], writes[0]}: 24,
	{writes[0], writes[6]}: 25,
	{writes[3], writes[6]}: 26,
	{writes[4], writes[1]}: 27,
	{writes[5], writes[5]}: 28,
	{writes[2], writes[2]}: 29,
	{writes[6], writes[1]}: 30,
	{writes[4], writes[3]}: 31,
	{writes[6], writes[0]}: 32,
	{writes[0], writes[3]}: 33,
	{writes[7], writes[2]}: 34,
	{writes[2], writes[7]}: 35,
	{writes[2], writes[4]}: 36,
	{writes[1], writes[2]}: 37,
	{writes[6], writes[5]}: 38,
	{writes[5], writes[3]}: 39,
	{writes[1], writes[6]}: 40,
	{writes[3], writes[4]}: 41,
	{writes[0], writes[1]}: 42,
	{writes[4], writes[0]}: 43,
	{writes[7], writes[1]}: 44,
	{writes[4], writes[7]}: 45,
	{writes[5], writes[1]}: 46,
	{writes[4], writes[5]}: 47,
	{writes[2], writes[1]}: 48,
	{writes[4], writes[2]}: 49,
	{writes[3], writes[3]}: 50,
	{writes[6], writes[6]}: 51,
	{writes[6], writes[4]}: 52,
	{writes[1], writes[3]}: 53,
	{writes[2], writes[3]}: 54,
	{writes[6], writes[2]}: 55,
	{writes[4], writes[4]}: 56,
	{writes[1], writes[1]}: 57,
	{writes[5], writes[4]}: 58,
	{writes[1], writes[5]}: 59,
	{writes[1], writes[4]}: 60,
	{writes[6], writes[3]}: 61,
	{writes[2], writes[5]}: 62,
	{writes[5], writes[2]}: 63,
}

func GetName(u uint8) string {
	if u == 0 {
		return "阳"
	} else {
		return "阴"
	}
}

func GetGuaName(u uint8) string {
	for k := range writes {
		if writes[k] == u {
			return bagua[k]
		}
	}
	return ""
}

func GetItemOFGua(p Pair) int {
	return GuaMap[p]*3
}
