package main

import (
  "fmt"
  "log"
  "encoding/csv"
  "os"
  "strings"
)

// remove index from slice by swapping element i with last index,
// then returning slice of everything but last index (what was at i) 
func remove(s [][]string, i int) [][]string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func prompt_for_file(msg string) (*os.File){
  fmt.Print(msg)
  var tname string
  _, lerr := fmt.Scanln(&tname)
  if lerr != nil {
    log.Fatal(lerr)
  }
  tfile, err := os.Open(tname)
  if err != nil {
    log.Fatal(err)
  }
  return tfile
}



func main() {
  type FoundAsset struct {
    location string
    full_serial string
  }
  type ColIndex struct {
    L_SER, A_SER, A_LOC int
  }
  // get lease list
  lease_file := prompt_for_file("Please enter the file path of the lease csv: ")
  asset_file := prompt_for_file("Please enter the file path of the asset csv: ")

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
  // "ASSETNUM":"MACH_SER_NUM" (full serial <-> full serial)
  // }
  // c2i = "columns to index"
  
  c2i := make(map[string]int,3)
  for k,v := range(assets[0]){
    if(v=="LOCATION" || v=="ASSETNUM" || v=="MACH_SER_NUM"){
      c2i[v]=k
    }
  }
  ci := ColIndex{5,0,33}

  // create map of len(lease list) to store locations
  loc_map := make(map[string]string, len(leases)-1)

  // for each serial in lease list, search for it in asset list
  // if lease serial HasSuffix asset serial:
  // add lease serial and asset location to map
  // remove index of asset list to avoid repetition in search


  fmt.Printf("There are %d assets in the lease list\n\n",len(leases))
  // for each serial in lease list, search asset list
  for _,lease_row := range(leases[1:]){
    var LR_SERIAL string
    LR_SERIAL = lease_row[ci.L_SER]
    // fmt.Printf("This asset has serial number %s\n",LR_SERIAL)
    
    for i:=1;i<len(assets);i++{
      // fmt.Printf("Lease: %s\tAsset %s\n",LR_SERIAL,assets[i][ci.A_SER])
      // if found, get LOCATION, store as string in map, continue
      // else mark "N/A" in map
      if strings.HasSuffix(assets[i][ci.A_SER],LR_SERIAL) {
        loc_map[LR_SERIAL]=assets[i][ci.A_LOC]
        // fmt.Printf("Location for asset %s set to %s\n",LR_SERIAL,assets[i][ci.A_LOC])
        remove(assets,i)
        continue
      } 
    }
        
    if val,ok := loc_map[LR_SERIAL]; !ok || strings.TrimSpace(val) == ""{
      // fmt.Printf("No value found at key %s\n",LR_SERIAL)
      loc_map[LR_SERIAL]="N/A"
    }
    
  }

  fmt.Printf("%d items in original lease listing\n", len(leases)-1)
  fmt.Printf("%d items in the location listing\n",len(loc_map))
  for k,v := range(loc_map){
    fmt.Printf("%s : %s\n",k,v)
  }

  // write new CSV with serial number and location
  // var conf_to_csv string
  // fmt.Println("Would you like to export to csv? [Y/n]")
  // _, err := fmt.Scanln(&conf_to_csv)
  // if conf_to_csv == "n" {
  //   log.Fatal("\nAll done!")
  // }
  
}