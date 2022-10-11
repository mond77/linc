# linc
a docker toy on Ubuntu(arm64) created for learning the principles of Docker runtime implementation.
use urface/cli sirupsen/logrus

To get codes of past versions see tags.

eg:process(implementation)

3.1           run-->init (Namespace)
details:
3.2,3.3       runProcess(CgroupManager,sendInitCmd),  Cgroups(memory,cpushare,cpuset)
details:
The 3rd part has implemented container process with Namespace and Cgroup.
