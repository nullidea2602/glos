local files = list_files()
print("Files in memory:")
for filename, _ in pairs(files) do
    print("-", filename)
end