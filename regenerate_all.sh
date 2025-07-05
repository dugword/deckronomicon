RUNS=100000

go run ./ -scenario thrill-of-possibility  -runs $RUNS
go run ./tools/analytics
mv results/summary.json thrill_of_possibility_summary.json

go run ./ -scenario stormshriek-feral -runs $RUNS
go run ./tools/analytics
mv results/summary.json flush_out_summary.json

go run ./ -scenario pearled-unicorn -runs $RUNS
go run ./tools/analytics
mv results/summary.json pearled_unicorn_summary.json
