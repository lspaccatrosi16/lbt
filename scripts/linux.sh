lbt
chmod +x ./out/lbt-linux_amd64
cp ./out/lbt-linux_amd64 ./lbt.temp
./lbt.temp
if [ $? -ne 0 ]; then
    exit 1
fi
cp out/lbt-linux_amd64 $(which lbt)
rm lbt.temp
