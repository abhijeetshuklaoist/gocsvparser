# gocsvparser

This is a tiny project to parse CSV files. [More details to come soon]

# Requiremet: 
1. To be able to parse CSVs with different kind of CSV having similar data. Some CSVs will have extra columns, some might be missing some required columns, some columns like Name might be split into multiple columns like FirstName and Last name
2. Read the CSV files, fetch Name, Email, Salary and Id
3. Dump records having above important fieled in one CSV and dump dirty data in another CSV


# Idea: 
The major issues which this problem poses are
1. Different set of columns: Different CSVs have different set of columns which might change, so we can't use static struct to map to different CSVs
2. Dirty data: We can't assume that input CSV is as we expected 
3. Scalability: Solution should be able to process large files 

1. Solution to first problem can be thought of as having a list  of all required columns and their possible different values. And a map of different possible values of same column name. 

Example: Name can be different columns FirstName and LastName, which itself can be in different format like fName or f.name or f_name. 
Name --> Name, name, firstname, lastname, f.name, l.name, first_name, last_name

When processing a file we can go through each header (actual header AH1) and get the expected header (EH1) from map. Now all we need to do is to read that record and parse that particular column as expected column EH1. 

In the case of, when there are not all expected column present we can just log the error and exit, we won't have to process the whole file in that case. Otherwise keep a slice in memory to store good and one slice to store bad data. Note: See if flush from slice can be done in parallel so that we don't face memory issue. 

2. Check for required fileds to be present in every line read from CSV and if not all required fields are present, log that line in a seprate file
3. The basic way we can scale this tool is by dividing the 3 major tasks among concurrent goroutines, i.e. reading the input file, writing correct data and writing incorrect data. 


# Performance number:
Input: A csv containing 2M records, with 20 columns, having 1M correct records, and 1M incorrect records. Size of file is arrprox 662MB
Expected output: 2 files, 1 having 1M correct records, another having 1M incorrect records
Execution time: On average (of 5 runs) on my machines, it took 12 seconds to process the file.
