v=v1.3.2
git tag $v
git push --tags
# go install 
go install github.com/ymzuiku/gojest@$v
echo "done."