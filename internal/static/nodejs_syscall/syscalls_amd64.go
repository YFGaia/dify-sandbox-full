//go:build linux && amd64

package nodejs_syscall

import "syscall"

const (
	//334
	SYS_RSEQ = 334
	// 435
	SYS_CLONE3 = 435
	// Additional syscalls for network operations
	SYS_GETRANDOM = 318
	SYS_SENDMMSG  = 307
)

var ALLOW_SYSCALLS = []int{
	// Basic file operations
	syscall.SYS_OPEN, syscall.SYS_WRITE, syscall.SYS_CLOSE, syscall.SYS_READ,
	syscall.SYS_OPENAT, syscall.SYS_NEWFSTATAT, syscall.SYS_IOCTL, syscall.SYS_LSEEK,
	syscall.SYS_FSTAT, syscall.SYS_STAT, syscall.SYS_LSTAT,
	syscall.SYS_GETDENTS64, syscall.SYS_READLINK, syscall.SYS_READLINKAT,
	syscall.SYS_ACCESS, syscall.SYS_FACCESSAT,

	// Memory management
	syscall.SYS_MPROTECT, syscall.SYS_MMAP, syscall.SYS_MUNMAP,
	syscall.SYS_MREMAP, syscall.SYS_MADVISE, syscall.SYS_BRK,

	// Signal handling
	syscall.SYS_RT_SIGACTION, syscall.SYS_RT_SIGPROCMASK,
	syscall.SYS_SIGALTSTACK, syscall.SYS_RT_SIGRETURN,

	// Process management
	syscall.SYS_GETPID, syscall.SYS_GETPPID, syscall.SYS_GETTID,
	syscall.SYS_GETUID, syscall.SYS_GETGID,
	syscall.SYS_SETUID, syscall.SYS_SETGID,
	syscall.SYS_EXIT, syscall.SYS_EXIT_GROUP,
	syscall.SYS_TGKILL, syscall.SYS_SCHED_YIELD,
	syscall.SYS_SCHED_GETAFFINITY,

	// Threading and synchronization
	syscall.SYS_FUTEX, syscall.SYS_SET_ROBUST_LIST, syscall.SYS_GET_ROBUST_LIST,
	SYS_RSEQ,

	// I/O multiplexing
	syscall.SYS_EPOLL_CREATE1, syscall.SYS_EPOLL_CTL, syscall.SYS_EPOLL_PWAIT,
	syscall.SYS_POLL, syscall.SYS_PPOLL, syscall.SYS_PSELECT6,

	// Time operations
	syscall.SYS_CLOCK_GETTIME, syscall.SYS_GETTIMEOFDAY, syscall.SYS_NANOSLEEP,
	syscall.SYS_CLOCK_NANOSLEEP, syscall.SYS_TIME,

	// File control
	syscall.SYS_FCNTL, syscall.SYS_DUP, syscall.SYS_DUP2, syscall.SYS_DUP3,
	syscall.SYS_PIPE, syscall.SYS_PIPE2,

	// Random number generation
	SYS_GETRANDOM,

	// Additional commonly needed syscalls for Node.js
	syscall.SYS_GETCWD, syscall.SYS_CHDIR,
	syscall.SYS_UMASK, syscall.SYS_GETRLIMIT, syscall.SYS_SETRLIMIT,
	syscall.SYS_GETRUSAGE, syscall.SYS_TIMES,
	syscall.SYS_UNAME, syscall.SYS_SYSINFO,

	// File system operations
	syscall.SYS_STATFS, syscall.SYS_FSTATFS,
	syscall.SYS_TRUNCATE, syscall.SYS_FTRUNCATE,
	syscall.SYS_FSYNC, syscall.SYS_FDATASYNC,

	// Additional syscalls that might be needed
	syscall.SYS_PRCTL, syscall.SYS_ARCH_PRCTL,
	syscall.SYS_GETPGRP, syscall.SYS_GETPGID, syscall.SYS_GETSID,

	// Allow all syscalls (same as Python) - complete freedom
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289, 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 315, 316, 317, 318, 319, 320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 346, 347, 348, 349, 350, 351, 352, 353, 354, 355, 356, 357, 358, 359, 360, 361, 362, 363, 364, 365, 366, 367, 368, 369, 370, 371, 372, 373, 374, 375, 376, 377, 378, 379, 380, 381, 382, 383, 384, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 397, 398, 399, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 419, 420, 421, 422, 423, 424, 425, 426, 427, 428, 429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 441, 442, 443, 444, 445, 446, 447, 448, 449, 450, 451, 452, 453, 454, 455, 456, 457, 458, 459, 460, 461, 462, 463, 464, 465, 466, 467, 468, 469, 470, 471, 472, 473, 474, 475, 476, 477, 478, 479, 480, 481, 482, 483, 484, 485, 486, 487, 488, 489, 490, 491, 492, 493, 494, 495, 496, 497, 498, 499, 500,
}

var ALLOW_ERROR_SYSCALLS = []int{
	syscall.SYS_CLONE, SYS_CLONE3,
	// Allow these to fail gracefully
	syscall.SYS_MKDIR, syscall.SYS_MKDIRAT,
	syscall.SYS_RMDIR, syscall.SYS_UNLINK, syscall.SYS_UNLINKAT,
	syscall.SYS_RENAME, syscall.SYS_RENAMEAT,
}

var ALLOW_NETWORK_SYSCALLS = []int{
	// Socket operations
	syscall.SYS_SOCKET, syscall.SYS_SOCKETPAIR,
	syscall.SYS_CONNECT, syscall.SYS_BIND, syscall.SYS_LISTEN,
	syscall.SYS_ACCEPT, syscall.SYS_ACCEPT4,

	// Data transfer
	syscall.SYS_SENDTO, syscall.SYS_RECVFROM,
	syscall.SYS_SENDMSG, syscall.SYS_RECVMSG,
	SYS_SENDMMSG,

	// Socket options and info
	syscall.SYS_GETSOCKNAME, syscall.SYS_GETPEERNAME,
	syscall.SYS_SETSOCKOPT, syscall.SYS_GETSOCKOPT,

	// Network utilities
	syscall.SYS_SHUTDOWN,
	syscall.SYS_UNAME,

	// File operations needed for network
	syscall.SYS_FCNTL, syscall.SYS_FSTATFS,
}
