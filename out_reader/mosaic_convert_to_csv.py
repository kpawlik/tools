import csv
import json
import os
import sys

CSV_COLUMNS = ["ASSOCIATED_NODE",
    "ICOMS_NUMBER",
    "STREET_ADDRESS",
    "MAC_ADDRESS",
    "DEVICE_TYPE",
    "DEVICE_DESCRIPTION",
    "POI_NUMBER",
    "TICKET_COUNT",
    "IMPACT_NUMBER",
    "VIOLATION_TYPE",
    "VIOLATION_DATE",
    "RULE_DESC",
    "RULE_NBR",
    "REMEDY",
    "POI_NODE_FAULT_ID",
    "POI_ELEMT_FAULT_ID",
    "SITE_ID",
    "FAULT_IMPACT_TYPE_ID",
    "POI_CNT",
    "GNIS_NODE_ID",
    "XMIT_INDICATOR",
    "POI_TYPE",
    "STATUS"]

def log(message):
    sys.stdout.write(message)

def out_file_name(raw_file_name):
    dir, name = os.path.split(raw_file_name)
    base_name = os.path.splitext(name)[0]
    csv_file_name = "{}.csv".format(base_name)
    return os.path.join(dir, csv_file_name)


def process_line(line):
    if not line:
        return 
    # print line
    line = line.strip()
    obj = json.loads(line)
    csv_line = []
    for column in CSV_COLUMNS:
        value = obj.get(column, '')
        if not value:
            value = obj.get(column.lower(), '')
        if column == "VIOLATION_DATE" and value.lower() == "none":
            value = ""
        if column == "FAULT_IMPACT_TYPE_ID":
            if value.lower() in  ["nsa", "non-service affecting"]:
                value = "NSA"
            if value.lower() in  ["sa", "service affecting"]:
                value = "SA"
            if value not in ["SA", "NSA"]:
                value = "UNC"
        csv_line.append(u"{}".format(value))
    return csv_line

def write_csv(file_name, csv_lines):
    with open(file_name, 'wb') as f:
        writer = csv.writer(f)
        writer.writerows(csv_lines)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Missing file path to convert")
        sys.exit(1)
    raw_file_name = sys.argv[1]
    csv_file_name = out_file_name(raw_file_name)
    csv_lines = []
    with open(raw_file_name, "r") as raw_file:
        for json_line in raw_file:
            try:
                csv_line = process_line(json_line)
                csv_lines.append(csv_line)
            except Exception as exc:

                log("Error parsing JSON: {}\n".format(exc))
                log("JSON line: {}\n".format(json_line))
    write_csv(csv_file_name, csv_lines)
    os.unlink(csv_file_name)
    os.unlink(raw_file_name)




