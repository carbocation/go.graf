paths=$(echo $GOPATH | tr ':' ' ')

for dir in static templates
do
	for path in $paths
	do
		fullpath="${path}/src/github.com/carbocation/go.gtfo/${dir}"
		if [ -d $fullpath ] && [ ! -d "./${dir}" ]
		then
			ln -s $fullpath ./
			continue
		fi
	done
done
