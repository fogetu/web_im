# install
go mod tidy 
go mod vendor
go build -mod=vendor

# run
./web_im

# 生成所有
bee generate scaffold team -fields="id:int64,tname:string,owner:string,create_at:int32" -driver=mysql -conn="debian-sys-maint:UTA5mehoU8fcvevB@tcp(127.0.0.1:3306)/web_im"

# 生成migration记录
bee migrate -driver=mysql -conn="debian-sys-maint:UTA5mehoU8fcvevB@tcp(127.0.0.1:3306)/web_im" -dir="database/migrations"


