
source /home/pfc15/Documents/aleatorio/go/spotgo-backend/venv/bin/activate
mkdir $1
spotdl --output $1  $2
cd $1
for f in *.mp3; do
    mv -- "$f" "$3.mp3"
done
cd ../../..
deactivate
