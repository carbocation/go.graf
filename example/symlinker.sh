paths=$(echo $GOPATH | tr ':' ' ')

for dir in static templates
do
	for path in $paths
	do
		fullpath="${path}/src/carbocation.com/code/go.asksite/${dir}"
		if [ -d $fullpath ] && [ ! -d "./${dir}" ]
		then
			ln -s $fullpath ./
			continue
		fi
	done
done
