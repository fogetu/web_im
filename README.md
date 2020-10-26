# install
go mod tidy 
go mod vendor
go build -mod=vendor
# run
./web_im

# migration
bee migrate -driver=mysql -conn="debian-sys-maint:UTA5mehoU8fcvevB@tcp(127.0.0.1:3306)/web_im" -dir="database/migrations"
