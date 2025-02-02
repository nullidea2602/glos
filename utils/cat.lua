-- Desc: Read a file and print its content

if not args[1] then
    print("Usage: lua print_file.lua <filename>")
else

    local file_to_read = args[1]
    local content, err = read_file(file_to_read)
    if content then
        print(content)
    else
        print("Error:", err)
    end

end