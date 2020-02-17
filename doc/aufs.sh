fs# pwd
/root/aufs
root@nicktming:~/aufs# ls 
run.sh
root@nicktming:~/aufs# cat run.sh
mkdir container-layer 
echo "I am container-layer" > container-layer/container-layer.txt

mkdir mnt

for i in {1..3}
do 
mkdir -p image-layer$i/subdir$i
echo "I am image layer $i" > image-layer$i/image-layer$i.txt
echo "subdir $i" > image-layer$i/subdir$i/subdir$i.txt
done
root@nicktming:~/aufs# ./run.sh 
root@nicktming:~/aufs# tree