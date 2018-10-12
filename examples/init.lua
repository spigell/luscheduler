local scenarios_dir = "examples/".."scenarios"

for _, a in pairs(filepath.glob(scenarios_dir.."/*.lua")) do

	local contents = ""
    local file = io.open( a, "r" )
    if (file) then
        contents = file:read()
        file:close()
        get_schedule = loadstring(contents)
        get_schedule()
    end
	schedule.new(SCHEDULE, a)
end

function sleep()
  os.execute("sleep infinity")
end
sleep()
