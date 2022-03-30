FROM  ubuntu

ADD . .

EXPOSE  2222
CMD ["./honeypot"]