package main

import (
  "fmt"
  "log"
  "encoding/csv"
  "os"
)

// remove index from slice by swapping element i with last index,
// then returning slice of everything but last index (what was at i) 
func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func main() {
  // get files
  
  // get lease list
  fmt.Print("Please enter the file path of the lease csv: ")
  var lease_name string
  _, lerr := fmt.Scanln(&lease_name)
  if lerr != nil {
    log.Fatal(lerr)
  }

  lease_file, err := os.Open(lease_name)
  if err != nil {
    log.Fatal(err)
  } 

  // get asset list
  fmt.Print("Please enter the file path of the asset csv: ")
  var asset_name string
  _, aerr := fmt.Scanln(&asset_name)
  if aerr != nil {
    log.Fatal(aerr)
  }

  asset_file, err := os.Open(asset_name)
  if err != nil {
    log.Fatal(err)
  }

  // read ALL of both files. Header at index 0
  lease_reader := csv.NewReader(lease_file)
  leases, err := lease_reader.ReadAll()
  if err != nil {
    log.Fatal(err)
  }
  asset_reader := csv.NewReader(asset_file)
  assets, err := asset_reader.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("The first column of the leases header is %s\n", leases[0][0])
  fmt.Printf("The first column of the assets header is %s\n", assets[0][0])

  // create map with known column headings to compare
  // Asset List:Lease List = {
  // "LOCATION"
  // "ASSETNUM":"Mfr. serial number" (full serial <-> full serial)
  // }

  // create map of len(lease list) to store locations

  // read all lines in from both 

  // for each serial in lease list, search asset list

  // if found, get LOCATION, store as string in map, continue,  

  // else mark "N/A" in map


}