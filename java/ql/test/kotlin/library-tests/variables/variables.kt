
class Foo {
    val prop: Int = 1

    fun myFunction(param: Int) {
        val local = 2 + 3
    }
}

val topLevel: Int = 1

class C1 {
    fun f1() {}
    fun f2() {}
}
class C2 (val o:C1) {
    fun f1() {}
    fun f3() {}

    fun f4() {
        o.ext();
    }
    fun C1.ext() {
        f1() // calls method defined in C1 class
        f2()
        f3()
        this.f1() // extensionReceiverParameter
        this.f2() // extensionReceiverParameter

        // calls method defined in C2 class
        this@C2.f1() // dispatchReceiverParameter
        this@C2.f3() // dispatchReceiverParameter
    }
}