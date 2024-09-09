rm -rf ./bin
make chaosd
make chaos-tools
mv bin chaosd-linux-amd64
tar -czvf chaosd-linux-amd64.tar.gz chaosd-v1.4.2-linux-amd64
protokitgo upload --source_files=chaosd-v1.4.2-linux-amd64.tar.gz --dest_dsn=oss@/chaosd/@pmt_setting@OSS_8 --config=setting.yaml --no_append_len --no_append_md5