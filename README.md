# Tugas Kecil 2 - Voxelization 3D dengan Octree

## Penjelasan Singkat
Program ini membaca file model 3D berformat OBJ, membangun struktur data Octree dari face segitiga, lalu menghasilkan representasi voxel (kumpulan kubus kecil) pada kedalaman tertentu.

Pendekatan utama yang digunakan adalah algoritma divide and conquer, yaitu membagi ruang 3D secara rekursif menjadi delapan sub-ruang (node anak Octree) hingga mencapai kedalaman yang ditentukan.

Pembentukan anak node Octree dijalankan secara concurrent (goroutine + WaitGroup) untuk meningkatkan performa proses.

Output program:
- File OBJ hasil voxelization di folder `test`.
- Statistik di terminal: jumlah voxel, vertex, faces, node terbentuk, node skipped, dan waktu eksekusi.

## Requirement dan Instalasi
Kebutuhan utama:
- Go (disarankan versi 1.20 atau lebih baru)

Struktur proyek yang digunakan:
- `src/main.go`
- `test/`

Langkah setup dependency:
1. Masuk ke root folder proyek.
2. Jika belum ada `go.mod`, inisialisasi modul:

```bash
go mod init tucil2
```

3. Rapikan dependency:

```bash
go mod tidy
```

## Cara Mengkompilasi Program
Jalankan dari root folder proyek:

```bash
go build -o bin/voxelizer.exe ./src/main.go
```

Setelah itu, file executable akan tersedia di folder `bin` dengan nama `voxelizer.exe`.

## Cara Menjalankan dan Menggunakan Program
### Menjalankan tanpa kompilasi
Jalankan dari root folder proyek:

```bash
go run ./src/main.go <path_file.obj> <depth>
```

Contoh:

```bash
go run ./src/main.go ./examples/model.obj 4
```

### Menjalankan hasil kompilasi

```bash
./bin/voxelizer.exe <path_file.obj> <depth>
```

Contoh:

```bash
./bin/voxelizer.exe ./examples/model.obj 4
```

### Format output file
Nama file output disimpan di folder `test` dengan format:

```text
<nama_asli>-voxelized-<depth>.obj
```

Contoh untuk `model.obj` dengan `depth` 4:

```text
test/model-voxelized-4.obj
```

## Author / Identitas Pembuat
- 18223003 - Wisa Ahmaduta Dinutama
- 18223066 - Nazwan Siddqi Muttaqin
