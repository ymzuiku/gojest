v=v1.1.19
git tag $v
git push --tags
go install github.com/ymzuiku/gojest@$v
echo "done."