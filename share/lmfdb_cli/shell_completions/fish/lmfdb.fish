# LMFDB CLI Fish Shell Completion

complete -c lmfdb -f -a "nf" -d "Query Number Fields"
complete -c lmfdb -f -a "ec" -d "Query Elliptic Curves"
complete -c lmfdb -f -a "query" -d "Generic query command"
complete -c lmfdb -f -a "list-collections" -d "List available API collections"
complete -c lmfdb -f -a "version" -d "Show version information"

# nf options
complete -c lmfdb -n "contains nf (commandline -op)" -l degree -s d -d "Number field degree" -r
complete -c lmfdb -n "contains nf (commandline -op)" -l discriminant -d "Filter by discriminant" -r
complete -c lmfdb -n "contains nf (commandline -op)" -l class-number -s h -d "Filter by class number" -r
complete -c lmfdb -n "contains nf (commandline -op)" -l limit -s n -d "Number of results" -r
complete -c lmfdb -n "contains nf (commandline -op)" -l fields -s f -d "Fields to return" -r
complete -c lmfdb -n "contains nf (commandline -op)" -l output -s o -d "Output file" -r -F
complete -c lmfdb -n "contains nf (commandline -op)" -l headless -d "Run browser in headless mode"
complete -c lmfdb -n "contains nf (commandline -op)" -l no-headless -d "Run browser in non-headless mode"

# ec options
complete -c lmfdb -n "contains ec (commandline -op)" -l rank -s r -d "Filter by rank" -r
complete -c lmfdb -n "contains ec (commandline -op)" -l torsion -s t -d "Filter by torsion" -r
complete -c lmfdb -n "contains ec (commandline -op)" -l conductor -d "Filter by conductor" -r
complete -c lmfdb -n "contains ec (commandline -op)" -l limit -s n -d "Number of results" -r
complete -c lmfdb -n "contains ec (commandline -op)" -l fields -s f -d "Fields to return" -r
complete -c lmfdb -n "contains ec (commandline -op)" -l output -s o -d "Output file" -r -F
complete -c lmfdb -n "contains ec (commandline -op)" -l headless -d "Run browser in headless mode"
complete -c lmfdb -n "contains ec (commandline -op)" -l no-headless -d "Run browser in non-headless mode"

# query options
complete -c lmfdb -n "contains query (commandline -op)" -l limit -s n -d "Number of results" -r
complete -c lmfdb -n "contains query (commandline -op)" -l fields -s f -d "Fields to return" -r
complete -c lmfdb -n "contains query (commandline -op)" -l output -s o -d "Output file" -r -F
complete -c lmfdb -n "contains query (commandline -op)" -l kwargs -s k -d "Additional query params as JSON" -r
complete -c lmfdb -n "contains query (commandline -op)" -l headless -d "Run browser in headless mode"
complete -c lmfdb -n "contains query (commandline -op)" -l no-headless -d "Run browser in non-headless mode"
