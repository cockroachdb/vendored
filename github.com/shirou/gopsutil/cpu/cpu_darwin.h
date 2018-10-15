typedef struct times_stat {
	int cpu;
	double user;
	double system;
	double idle;
	double nice;
} times_stat;

times_stat* per_cpu_times(int *out_ncpus);
times_stat all_cpu_times(void);
