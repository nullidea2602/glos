-- lwrite reads input line by line until ":exit" then writes to filename from args[1]
if not args or #args == 0 then
    print("Usage: lwrite <filename>")
    return
end

local filename = args[1]
local content = ""

print("Enter content line by line. Type ':exit' to finish.")
content = read_multiline_input()

write_file(filename, content)