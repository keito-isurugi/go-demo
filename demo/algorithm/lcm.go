package algorithm

import "fmt"

// LCMDemo calculates the Least Common Multiple of two numbers
// using the formula: LCM(a, b) = (a * b) / GCD(a, b)
func LCMDemo(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	// Make positive for calculation
	absA, absB := absHelper(a), absHelper(b)
	return (absA * absB) / GCD(absA, absB)
}

// LCMMultipleDemo calculates the LCM of multiple numbers
func LCMMultipleDemo(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}
	if len(numbers) == 1 {
		return absHelper(numbers[0])
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = LCMDemo(result, numbers[i])
	}
	return result
}

// LCMByPrimeFactorization calculates LCM using prime factorization
func LCMByPrimeFactorization(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}

	absA, absB := absHelper(a), absHelper(b)
	factorsA := primeFactors(absA)
	factorsB := primeFactors(absB)

	// Merge factors taking maximum power for each prime
	maxFactors := make(map[int]int)
	for prime, count := range factorsA {
		maxFactors[prime] = count
	}
	for prime, count := range factorsB {
		if maxFactors[prime] < count {
			maxFactors[prime] = count
		}
	}

	// Calculate LCM from factors
	result := 1
	for prime, count := range maxFactors {
		for i := 0; i < count; i++ {
			result *= prime
		}
	}
	return result
}

// primeFactors returns a map of prime factors and their counts
func primeFactors(n int) map[int]int {
	factors := make(map[int]int)

	// Check for factor 2
	for n%2 == 0 {
		factors[2]++
		n = n / 2
	}

	// Check for odd factors
	for i := 3; i*i <= n; i += 2 {
		for n%i == 0 {
			factors[i]++
			n = n / i
		}
	}

	// If n is a prime greater than 2
	if n > 2 {
		factors[n]++
	}

	return factors
}

// absHelper returns absolute value
func absHelper(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// DemoLCM demonstrates various LCM calculations
func DemoLCM() {
	fmt.Println("=== LCM (Least Common Multiple) Demo ===")
	fmt.Println()

	// Basic LCM examples
	fmt.Println("Basic LCM Examples:")
	pairs := [][2]int{
		{12, 18},
		{15, 20},
		{7, 5},
		{24, 36},
		{10, 15},
	}

	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		lcm := LCMDemo(a, b)
		gcd := GCD(a, b)
		fmt.Printf("LCM(%d, %d) = %d (GCD = %d)\n", a, b, lcm, gcd)
		fmt.Printf("  Verification: %d × %d = %d × %d = %d\n",
			lcm/a, a, lcm/b, b, lcm)
	}

	fmt.Println("\nLCM using Prime Factorization:")
	for _, pair := range pairs[:3] {
		a, b := pair[0], pair[1]
		lcm := LCMByPrimeFactorization(a, b)
		fmt.Printf("LCM(%d, %d) = %d (using prime factorization)\n", a, b, lcm)
	}

	fmt.Println("\nLCM of Multiple Numbers:")
	multiSets := [][]int{
		{4, 6, 8},
		{3, 5, 7},
		{12, 15, 20},
		{2, 3, 4, 5},
	}

	for _, set := range multiSets {
		lcm := LCMMultipleDemo(set...)
		fmt.Printf("LCM%v = %d\n", set, lcm)
	}

	fmt.Println("\nPractical Applications:")

	// Application 1: Finding common period
	fmt.Println("\n1. Finding Common Period:")
	fmt.Println("   Train A arrives every 12 minutes")
	fmt.Println("   Train B arrives every 18 minutes")
	lcm := LCMDemo(12, 18)
	fmt.Printf("   They meet at the platform every %d minutes\n", lcm)

	// Application 2: Gear ratios
	fmt.Println("\n2. Gear System:")
	fmt.Println("   Gear A has 15 teeth")
	fmt.Println("   Gear B has 20 teeth")
	lcm = LCMDemo(15, 20)
	fmt.Printf("   After %d teeth pass the contact point,\n", lcm)
	fmt.Printf("   Gear A completes %d rotations\n", lcm/15)
	fmt.Printf("   Gear B completes %d rotations\n", lcm/20)

	// Application 3: LCD fractions
	fmt.Println("\n3. Adding Fractions (finding LCD):")
	fmt.Println("   1/6 + 1/9 = ?")
	lcd := LCMDemo(6, 9)
	fmt.Printf("   LCD = %d\n", lcd)
	fmt.Printf("   1/6 + 1/9 = %d/%d + %d/%d = %d/%d\n",
		lcd/6, lcd, lcd/9, lcd, (lcd/6)+(lcd/9), lcd)
}