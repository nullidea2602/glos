-- rm.lua deletes file from args[1] using the delete_file("filename") api

if not args or #args == 0 then
    print("Usage: rm <filename>")
    return
end

local file_to_delete = args[1]
local success, err = delete_file(file_to_delete)

if not success then
    print("Error:", err or "Unknown error") -- Handle nil case explicitly
    return
end
