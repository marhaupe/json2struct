package cmd

// this is not right; TODO!
// type JSONToStruct []struct {
// 	In_array struct {
// 		Input_index            int    `json:"input_index"`
// 		Candidate_index        int    `json:"candidate_index"`
// 		Delivery_line_1        string `json:"delivery_line_1"`
// 		Last_line              string `json:"last_line"`
// 		Delivery_point_barcode string `json:"delivery_point_barcode"`
// 		Components             struct {
// 			Primary_number             string `json:"primary_number"`
// 			Street_predirection        string `json:"street_predirection"`
// 			Street_name                string `json:"street_name"`
// 			Street_suffix              string `json:"street_suffix"`
// 			City_name                  string `json:"city_name"`
// 			State_abbreviation         string `json:"state_abbreviation"`
// 			Zipcode                    string `json:"zipcode"`
// 			Plus4_code                 string `json:"plus4_code"`
// 			Delivery_point             string `json:"delivery_point"`
// 			Delivery_point_check_digit string `json:"delivery_point_check_digit"`
// 		} `json:"components"`
// 		Metadata struct {
// 			Record_type            string `json:"record_type"`
// 			Zip_type               string `json:"zip_type"`
// 			County_fips            string `json:"county_fips"`
// 			County_name            string `json:"county_name"`
// 			Carrier_route          string `json:"carrier_route"`
// 			Congressional_district string `json:"congressional_district"`
// 			Rdi                    string `json:"rdi"`
// 			Elot_sequence          string `json:"elot_sequence"`
// 			Elot_sort              string `json:"elot_sort"`
// 			Latitude               int    `json:"latitude"`
// 			Longitude              int    `json:"longitude"`
// 			Precision              string `json:"precision"`
// 			Time_zone              string `json:"time_zone"`
// 			Utc_offset             int    `json:"utc_offset"`
// 			Dst                    bool   `json:"dst"`
// 		} `json:"metadata"`
// 		Analysis struct {
// 			Dpv_match_code string `json:"dpv_match_code"`
// 			Dpv_footnotes  string `json:"dpv_footnotes"`
// 			Dpv_cmra       string `json:"dpv_cmra"`
// 			Dpv_vacant     string `json:"dpv_vacant"`
// 			Active         string `json:"active"`
// 		} `json:"analysis"`
// 	} `json:"in_array,omitempty"`
// 	In_array struct {
// 		Input_index            int    `json:"input_index"`
// 		Candidate_index        int    `json:"candidate_index"`
// 		Delivery_line_1        string `json:"delivery_line_1"`
// 		Last_line              string `json:"last_line"`
// 		Delivery_point_barcode string `json:"delivery_point_barcode"`
// 		Components             struct {
// 			Primary_number             string `json:"primary_number"`
// 			Street_predirection        string `json:"street_predirection"`
// 			Street_name                string `json:"street_name"`
// 			Street_suffix              string `json:"street_suffix"`
// 			City_name                  string `json:"city_name"`
// 			State_abbreviation         string `json:"state_abbreviation"`
// 			Zipcode                    string `json:"zipcode"`
// 			Plus4_code                 string `json:"plus4_code"`
// 			Delivery_point             string `json:"delivery_point"`
// 			Delivery_point_check_digit string `json:"delivery_point_check_digit"`
// 		} `json:"components"`
// 		Metadata struct {
// 			Record_type            string `json:"record_type"`
// 			Zip_type               string `json:"zip_type"`
// 			County_fips            string `json:"county_fips"`
// 			County_name            string `json:"county_name"`
// 			Carrier_route          string `json:"carrier_route"`
// 			Congressional_district string `json:"congressional_district"`
// 			Rdi                    string `json:"rdi"`
// 			Elot_sequence          string `json:"elot_sequence"`
// 			Elot_sort              string `json:"elot_sort"`
// 			Latitude               int    `json:"latitude"`
// 			Longitude              int    `json:"longitude"`
// 			Precision              string `json:"precision"`
// 			Time_zone              string `json:"time_zone"`
// 			Utc_offset             int    `json:"utc_offset"`
// 			Dst                    bool   `json:"dst"`
// 		} `json:"metadata"`
// 		Analysis struct {
// 			Dpv_match_code string `json:"dpv_match_code"`
// 			Dpv_footnotes  string `json:"dpv_footnotes"`
// 			Dpv_cmra       string `json:"dpv_cmra"`
// 			Dpv_vacant     string `json:"dpv_vacant"`
// 			Active         string `json:"active"`
// 		} `json:"analysis"`
// 	} `json:"in_array,omitempty"`
// }
