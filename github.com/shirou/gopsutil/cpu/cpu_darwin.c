#include <stdlib.h>
#include <mach/mach_host.h>
#include <mach/host_info.h>
#include <libproc.h>
#include <mach/processor_info.h>
#include <mach/vm_map.h>
#include "cpu_darwin.h"

times_stat* per_cpu_times(int *out_ncpus) {
	natural_t ncpus;
	processor_info_array_t info;
	mach_msg_type_number_t infosz;

	kern_return_t status = host_processor_info(mach_host_self(),
		PROCESSOR_CPU_LOAD_INFO, &ncpus, &info, &infosz);
	if (status != KERN_SUCCESS) {
		return NULL;
	}
	*out_ncpus = ncpus;

	times_stat* out = malloc(sizeof(times_stat) * ncpus);
	if (out == NULL) {
		goto done;
	}

	processor_cpu_load_info_data_t* cpuloads = (processor_cpu_load_info_data_t*) info;
	natural_t i;
	for (i = 0; i < ncpus; i++) {
		out[i].cpu = i;
		out[i].user = ((double) cpuloads[i].cpu_ticks[CPU_STATE_USER]) / CLOCKS_PER_SEC;
		out[i].system = ((double) cpuloads[i].cpu_ticks[CPU_STATE_SYSTEM]) / CLOCKS_PER_SEC;
		out[i].idle = ((double) cpuloads[i].cpu_ticks[CPU_STATE_IDLE]) / CLOCKS_PER_SEC;
		out[i].nice = ((double) cpuloads[i].cpu_ticks[CPU_STATE_NICE]) / CLOCKS_PER_SEC;
	}

done:
	vm_deallocate(mach_task_self(), (vm_address_t) info, infosz);
	return out;
}

times_stat all_cpu_times(void) {
	host_cpu_load_info_data_t cpuload;
	mach_msg_type_number_t cpuloadsz;

	kern_return_t status = host_statistics(mach_host_self(), HOST_CPU_LOAD_INFO,
		(host_info_t) &cpuload, &cpuloadsz);
	if (status != KERN_SUCCESS) {
		return (times_stat) {};
	}

	times_stat out;
	out.cpu = 0;
	out.user = ((double) cpuload.cpu_ticks[CPU_STATE_USER]) / CLOCKS_PER_SEC;
	out.system = ((double) cpuload.cpu_ticks[CPU_STATE_SYSTEM]) / CLOCKS_PER_SEC;
	out.idle = ((double) cpuload.cpu_ticks[CPU_STATE_IDLE]) / CLOCKS_PER_SEC;
	out.nice = ((double) cpuload.cpu_ticks[CPU_STATE_NICE]) / CLOCKS_PER_SEC;
	return out;
}
