use 'test'

create table ticks
(
	id int auto_increment,
	ticks timestamp not null,
	symbol enum('ETH-BTC', 'BTC-USD', 'BTC-EUR') not null,
	bid float(64,2) not null,
	`ask(64)` float(64,2) not null,
	constraint ticks_pk
		primary key (id)
);