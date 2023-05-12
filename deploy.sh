v=v1.2.1
git tag $v
git push --tags
# go install 
go install github.com/ymzuiku/gojest@$v
echo "done."