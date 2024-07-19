FROM scratch
ADD bin/leaderboard_amd64 /leaderboard
CMD ["/leaderboard"]
