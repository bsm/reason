#!/usr/bin/env ruby

require 'csv'

RND = Random.new(100)

def x(mean, stddev = mean/5)
	if @last
		y1 = @last
		@last = nil
	else
		w = 1
		until w < 1.0 do
			x1 = 2.0 * RND.rand - 1.0
			x2 = 2.0 * RND.rand - 1.0
			w  = x1 * x1 + x2 * x2
		end
		w = Math.sqrt((-2.0 * Math.log(w)) / w)
		y1 = x1 * w
		@last = x2 * w
	end
	return mean + y1 * stddev
end

OUT = CSV.new(STDOUT)
CSV.foreach(ARGV[0]) do |row|
	c1, c2, c3, c4, c5, n1, n2, n3, n4, n5 = \
		row[0, 10].map {|s| s.tr('v', '').to_f }
	num = c1*x(0.01) +
				c2*x(2.6, 0.2) +
				c3*x(0.1) +
				c4*x(0.8) +
				c5*x(0.02) +
				(n1+1)*n4*x(2, 0.1) +
				n2*x(0.2) +
				n3*x(0.1) +
				n5*x(0.3)
	OUT << row + [num.round(1)]
end
