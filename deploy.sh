v=v1.1.16
git tag $v
git push --tags
go install github.com/ymzuiku/gojest@$v
echo "done."