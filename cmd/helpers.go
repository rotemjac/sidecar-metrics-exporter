package main

import (
    "strings"
)

func mergeToOneList(stringGroups ...string) []string {
    var mergedList []string
    for _, str := range stringGroups {
        splitItems := strings.Split(str, ";")
        for _, item := range splitItems {
            if item != "" {
                mergedList = append(mergedList, item)
            }
        }
    }
    return mergedList
}
