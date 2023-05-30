use wapc_guest as wapc;

#[no_mangle]
pub fn wapc_init() {
  wapc::register_function("hello", hello);
}

//Metodo generico donde se lee el mensaje obtenido
fn hello(msg: &[u8]) -> wapc::CallResult {
  wapc::console_log(&format!(
    "Se recibe: {}",
    std::str::from_utf8(msg).unwrap()
  ));

  let mut x = 0;

  //Procedimiento para convertir un mensaje en bytes a un array de strings sin espacios para poder identificar cada palabra
  let s = Vec::from(msg);
  let snew = String::from_utf8(s).unwrap();
  let v: Vec<&str> = snew.split(' ').collect();

  //Iteramos el array de String en busca de palabras clave
  for i in v.iter() {
    match i {
        &"convertir" => x = 1,
        &"conectado" => x = 2,
        &"conectados" => x = 2,
        &"tarea" => x = 3,
        &"tareas" => x = 3,
        &"encender" => x= 4,
        &"pausar" => x= 4,
        &"apagar" => x= 4,
        &"estado" => x= 4,
        &"reproducir"=> x=4,
        &"haz" => x = 5,
        &"conectes" => x = 5,
        &&_ => println!("Hecho"),
    }
  }

  //Situamos x = 1 con el metodo para cambiar de divisa
  if x == 1 {
    let r: f32 = convertir(msg);
    let r1  = r.to_string();
    let r2 = r1.as_bytes();
    let _res = wapc::host_call("binding", "sample:namespace", "pong", r2)?;
    Ok(r2.to_vec())

  } else if x == 2{
    //Se quiere diferenciar entre si se pregunta por alguien en particular o la lista de todos los conectados
    let mut y = 0;
    for j in v.iter(){
      match j {
        &"quien" => y = 1,
        &"lista" => y = 1,
        &&_ => println!("Hecho")
      }
    }
    if y == 1 {
      //En este caso sera una pregunta del tipo quien esta conectado o lusta de conectados
      let conectados = String::from("conectado");
      let byteconectados = conectados.as_bytes();
      Ok(byteconectados.to_vec())
    } else {
      //Aqui se habra preguntado por ESTA X CONECTADO entonces devolveremos X CONECTADO
      let s = Vec::from(msg);
      let mut snew = String::from_utf8(s).unwrap();
      let snew1 = snew.split_off(5);
      let snew2 = snew1.as_bytes();
      Ok(snew2.to_vec())
    }
  } else if x == 3{
    let mut z = 0;
    for j in v.iter(){
      match j {
        &"eliminar" => z = 1,
        &&_ => println!("Hecho")
      }}
      if z == 1 {
        let eliminartarea = String::from("eliminartarea");
        let byteeliminartarea = eliminartarea.as_bytes();
        Ok(byteeliminartarea.to_vec())
      } else {
        let tareas = String::from("tareas");
        let bytetareas = tareas.as_bytes();
        Ok(bytetareas.to_vec())
      }

  } else if x == 4 {
    let accion = String::from("accion");
    let byteaccion = accion.as_bytes();
    Ok(byteaccion.to_vec())
  } else if x ==5  {
    let enviartarea = String::from("enviartarea");
    let byteenviartarea = enviartarea.as_bytes();
    Ok(byteenviartarea.to_vec())
  } else {
    //Si no identificamos ninguna palabra clave se devolvera que no se entendió la pregunta
      let string = String::from("no entiendo la pregunta");
      let u8s = string.as_bytes();
      let _res = wapc::host_call("binding", "sample:namespace", "pong", u8s)?;
      Ok(u8s.to_vec())
    }
  
  
  
}

  //Procedimiento para realizar el cambio de divisa
  fn convertir (msg: &[u8]) -> f32 {

    //Procedimiento para convertir un mensaje en bytes a un array de strings sin espacios para poder identificar cada palabra
    let s = Vec::from(msg);
    let mut snew = String::from_utf8(s).unwrap();
    println!("{}", snew);
    let snew1 = snew.split_off(10);
    let v: Vec<&str> = snew1.split(' ').collect();


    //Obtenemos el numero y las divisas del mensaje que deberá ser predefinido
    let numero: &str = v[0];
    let numeroc = String::from_utf8(numero.into()).unwrap();
    let my_int = numeroc.parse::<f32>().unwrap();
    
    let divisa1: &str = v[1];
    let sdivisa1 = String::from_utf8(divisa1.into()).unwrap();
    let divisa2: &str = v[3];
    let sdivisa2 = String::from_utf8(divisa2.into()).unwrap();
 
    let sdividendo = sdivisa2;
    let sdivisor = sdivisa1;
  
    //Realizamos el cálculo del cambio de divisa en función de la moneda que se pida
    let mut dividendo: f32 = 0.0;
    let mut divisor: f32 = 1.0;
    
    if sdividendo == String::from("euros"){
      dividendo = 1 as f32;
    }
    if sdividendo == String::from("dolares"){
      dividendo = 1.06 as f32;
    }
    if sdividendo == String::from("libras"){
      dividendo = 0.88 as f32;
    }
    if sdividendo == String::from("yenes"){
      dividendo = 141.29 as f32;
    }
    if sdivisor == String::from("euros"){
      divisor = 1 as f32;
    }
    if sdivisor == String::from("dolares"){
      divisor = 1.06 as f32;
    }
    if sdivisor == String::from("libras"){
      divisor = 0.88 as f32;
    }
    if sdivisor == String::from("yenes"){
      divisor = 141.29 as f32;
    }
  
    let resultado: f32 = (my_int*dividendo)/divisor;
  
    println!("{}", resultado);
    return resultado;
     

  }