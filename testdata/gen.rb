require 'csv'
require 'zlib'

rnd = Random.new(100)
pos = 0

OUT = CSV.new(File.open('bigdata.csv', 'w'))
CSV.foreach('bigcls.csv') do |row|
	num = row[10] == 'c1' ? 1.0 : 3.0

	row[0,5].each_with_index do |s, i|
		n = s.tr('v', '').to_f
		num += rnd.rand.fdiv(n+0.1)
	end

	row[5,5].each_with_index do |s, i|
		n = s.to_f
		num += n #  * (rnd.rand+0.5).fdiv(50-i))
	end

	OUT << row + [num.round(2)]
	pos += 1

	# break if pos > 100
end
