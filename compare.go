package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// remove index from slice by swapping element i with last index,
// then returning slice of everything but last index (what was at i)
func remove(s [][]string, i int) [][]string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func prompt_for_input(msg string) *os.File {
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

func prompt_for_output(msg string) *os.File {
	fmt.Print(msg)
	var tname string
	_, lerr := fmt.Scanln(&tname)
	if lerr != nil {
		log.Fatal(lerr)
	}
	tfile, err := os.Create(tname)
	if err != nil {
		log.Fatal(err)
	}
	return tfile
}

func main() {

	type FoundAsset struct {
		location    string
		full_serial string
	}
	type ColIndex struct {
		L_SER, A_SER, A_LOC int
	}

	lease_file := prompt_for_input("Please enter the file path of the lease csv: ")
	asset_file := prompt_for_input("Please enter the file path of the asset csv: ")

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

	ci := ColIndex{5, 0, 33}

	// store locations to map and a list of full details
	loc_map := make(map[string]string, len(leases)-1)

	// create header and fill
	new_header := make([]string, len(assets[0])+len(leases[0]), len(assets[0])+len(leases[0]))
	for i, h := range leases[0] {
		new_header[i] = h
	}
	for i, h := range assets[0] {
		new_header[i+len(leases[0])] = h
	}

	// create full details list
	full_details := make([][]string, 0, len(leases)-1)

	fmt.Printf("There are %d assets in the lease list\n\n", len(leases))
	// for each serial in lease list, search asset list

	for _, lease_row := range leases[1:] {
		var LR_SERIAL string
		LR_SERIAL = lease_row[ci.L_SER]

		for i := 1; i < len(assets); i++ {
			if strings.HasSuffix(assets[i][ci.A_SER], LR_SERIAL) {
				loc_map[LR_SERIAL] = assets[i][ci.A_LOC]

				// full_detail data
				fdt := make([]string, len(assets[0])+len(leases[0]), len(assets[0])+len(leases[0]))
				for l := 0; l < len(lease_row); l++ {
					tlr := lease_row[l]
					fdt[l] = tlr
				}
				for l := 0; l < len(assets[i]); l++ {
					ail := assets[i][l]
					fdt[l+len(lease_row)] = ail
				}
				full_details = append(full_details, fdt)
				remove(assets, i)
				continue
			}
		}

		if val, ok := loc_map[LR_SERIAL]; !ok || strings.TrimSpace(val) == "" {
			loc_map[LR_SERIAL] = "N/A"
			full_details = append(full_details, []string{"N/A"})

		}
	}

	fmt.Printf("%d items in original lease listing\n", len(leases)-1)
	fmt.Printf("%d items in the location listing\n", len(loc_map))
	for k, v := range loc_map {
		fmt.Printf("%s : %s\n", k, v)
	}

	// write new CSV with serial number and location
	var conf_to_csv string
	fmt.Println("Would you like to export to csv? [Y/n]")
	_, cerr := fmt.Scanln(&conf_to_csv)
	if conf_to_csv == "n" || cerr != nil {
		log.Fatal("\nAll done!")
	}

	output_file := prompt_for_output("Where would you like to save this list? ")
	csv_out := csv.NewWriter(output_file)
	var csv_type string
	fmt.Println("[1] Asset and location\n[2]Full details")
	_, cerr = fmt.Scanln(&csv_type)
	if csv_type == "1" {
		csv_out.Write([]string{"Asset", "Location"})
		new_records := make([][]string, 0, len(loc_map))
		for k, v := range loc_map {
			new_records = append(new_records, []string{k, v})
		}
		csv_out.WriteAll(new_records)
		csv_out.Flush()
	} else if csv_type == "2" {
		csv_out.Write(new_header)
		new_records := make([][]string, 0, len(full_details))
		for _, v := range full_details {
			new_records = append(new_records, v)
		}
		csv_out.WriteAll(new_records)
		csv_out.Flush()
	} else {
		log.Fatal("Issue parsing reqeuested data type")
	}
}
