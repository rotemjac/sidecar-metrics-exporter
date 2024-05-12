package main

func mergeLists(listsGroups ...[]string) []string {
    var mergedList []string
    for _, list := range listsGroups {
        for _, item := range list {
            splitItems := strings.Split(item, ";")
            mergedList = append(mergedList, splitItems...)
        }
    }
    return mergedList
}
