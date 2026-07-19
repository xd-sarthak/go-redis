package core

// supports 4 databases
var KeySpaceStat [4]map[string]int

func UpdateDBStats(num int, metric string, value int) {
	KeySpaceStat[num][metric] = value
}