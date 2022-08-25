CREATE TABLE IF NOT EXISTS board(
   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
   size_x INTEGER,
   size_y INTEGER
);

CREATE TABLE IF NOT EXISTS robot(
   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
   x INTEGER,
   y INTEGER,
	direction TEXT,
	board_id UUID REFERENCES board (id)
);
