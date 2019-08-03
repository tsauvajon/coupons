INSERT INTO brands(name)
vALUES  ('Tesco'),
        ('Morrisons'),
        ('Sainsbury'),
        ('ALDI'),
        ('Asda'),
        ('LIDL'),
        ('Iceland');

INSERT INTO coupons(name, value, expiry_utc, brand_id)
VALUES  ('Save £10 at Tesco',  10, '2019-09-25 08:01:13', (SELECT id FROM brands WHERE name = 'Tesco')),
        ('Save £20 at Tesco',  20, '2019-03-01 10:15:53', (SELECT id FROM brands WHERE name = 'Tesco')),
        ('Save £30 at Tesco',  30, '2019-11-11 10:15:53', (SELECT id FROM brands WHERE name = 'Tesco')),
        ('Save £1 at LIDL',  1, '2019-09-25 10:15:53', (SELECT id FROM brands WHERE name = 'LIDL')),
        ('Save £2 at LIDL',  2, '2019-10-02 10:15:53', (SELECT id FROM brands WHERE name = 'LIDL')),
        ('Save £3 at LIDL',  3, '2019-11-11 10:15:53', (SELECT id FROM brands WHERE name = 'LIDL')),
        ('Save £4 at LIDL',  4, '2019-09-25 08:01:13', (SELECT id FROM brands WHERE name = 'LIDL')),
        ('Save £5 at LIDL',  5, '2019-10-02 08:01:13', (SELECT id FROM brands WHERE name = 'LIDL')),
        ('Save £30 at Asda',  30, '2019-11-11 08:01:13', (SELECT id FROM brands WHERE name = 'Asda'));
