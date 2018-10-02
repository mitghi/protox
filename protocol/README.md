# protocol

# PDU ( protocol data unit )

## Queue
#### structure

	MessageId	: uint16 
	Address		: string
	ReturnPath	: string
	Mark		: bytes
	Message		: bytes


### fixed header

+ fixed header **options**:

	| bit | semantic |
	|:---|---:|
	|**0**| |		
	|**1**| hasPayload |	
	|**2**| isDuplicate |	
	|**3**| hasOpts |		

### variable header

Options: **1** byte
  
+  *options* (0xF0):
	        
	| bit | semantic |
	|:---|---:|
	|**0**| hasMark |
	|**1**| hasReturnPath |
	|**2**| hasAddress |
	|**3**| hasId |

+ *command* (0x0F):

	| bit | semantic |
	|:---|---:|
	|**4**| NOP |
	|**5**| Drain |
	|**6**| Destroy |
	|**7**| Initialize |

### types

+ payload **item** types:

	|payload|type|
	|:---|---:|
	|**MessageId**| unsigned int32|
	|**Address**| string|
	|**ReturnPath**| string|
	|**Mark**| bytes|
	|**Message**|bytes|


