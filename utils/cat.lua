if not args or #args == 0 then
    print("Usage: lua print_file.lua <filename>")
    return
end

local file_to_read = args[1]
local content, err = read_file(file_to_read)

if not content then
    print("Error:", err)
    return
end

print(content)